import React, {useEffect, useState} from 'react';

import {
  Divider,
  AppBar,
  Toolbar,
  Button,
  Drawer,
  TextField,
  FormControl,
  FormHelperText,
  IconButton,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  MenuItem,
  Select,
  Table,
  Typography,
  TableRow,
  TableCell,
  TableHead,
  TableBody,
  TableContainer,
} from '@material-ui/core';

import {Sort, Menu, ChevronLeft} from '@material-ui/icons';

import './App.css';


const SortVisualizer = (props) => {
  const {steps} = props;
  return (
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow key='header'>
            <TableCell colSpan={steps[0].list.length}>Steps</TableCell>
          </TableRow>

        </TableHead>
        <TableBody>
          {steps.map((step) => (
            <TableRow key={step.id}>
              {step.list.map((value) => (
                <TableCell key={step.id + value}>{value}</TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};


/**
 * Sends a POST request to the server to get sorting steps.
 * @param {string} input  a comma-delimited string containing numbers
 * @param {string} algorithm  algorithm to use for sorting the numbers
 * @return {Promise}  the array with the step by step sorting
 */
function requestSort(input, algorithm) {
  return new Promise((accept, reject) => {
    const request = new XMLHttpRequest();
    request.onload = function() {
      const stepData = JSON.parse(request.responseText);
      accept(stepData['steps']);
    };
    request.onerror = function(err) {
      reject(err);
    };
    request.open('POST', 'http://127.0.0.1:8080/sort');
    const isNotNumber = (element) => element !== '' || !isNaN(element);
    const elements = input.split(',').filter(isNotNumber);
    const sortItems = new Set(elements.map((element) => parseFloat(element)));
    request.send(JSON.stringify({
      elements: Array.from(sortItems),
      algorithm,
    }));
  });
}

/**
 * Returns an array of supported sorting algorithms.
 * @return {Promise}
 */
function requestSupported() {
  return new Promise((accept, reject) => {
    const request = new XMLHttpRequest();
    request.onload = function() {
      const supportedList = JSON.parse(request.responseText);
      accept(supportedList);
    };
    request.onerror = function(err) {
      reject(err);
    };
    request.open('GET', 'http://127.0.0.1:8080/sort/supported');
    request.send();
  });
}

const SortingPlugin = (props) => {
  const [algorithm, setAlgorithm] = useState(props.algorithm);
  const [steps, setSteps] = useState(props.steps || []);
  const [supported, setSupported] = useState([]);

  useEffect(() => {
    // ask server for supported algorithms
    requestSupported().then((supportedList) => {
      setSupported(supportedList);
    });
  }, [supported]);

  const [input, setInput] = useState('');

  const handleClick = () => {
    // request sorting operation from server
    requestSort(input, algorithm)
        .then((newSteps) => setSteps(newSteps))
        .catch((err) => console.error(err));
  };

  return (
    <div className='center-screen'>
      <TextField
        value={input}
        onChange={(event) => setInput(event.target.value)}
        multiline
        minRows={4}
        label='Input comma separated integers'
        id='input-elements'
      />
      <br/>
      <FormControl>
        <Select
          labelId='algorithm-select-label'
          id="algorithm-select"
          value={algorithm}
          onChange={(event) => setAlgorithm(event.target.value)}
        >
          {supported.map((algo) => {
            return <MenuItem key={algo} value={algo}>{algo}</MenuItem>;
          })}
        </Select>
        <FormHelperText>Sorting Algorithm</FormHelperText>
      </FormControl>
      <br/>
      <Button onClick={handleClick}>Submit</Button>
      {steps.length ? <SortVisualizer steps={steps}/> : null}
    </div>
  );
};

const App = () => {
  // ask server for starting algorithm
  const container = window.document.body || undefined;

  const [openDrawer, setDrawerOpen] = useState(false);
  return (
    <div>
      <AppBar
        position="fixed"
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            onClick={() => setDrawerOpen(!openDrawer)}
            edge="start"
          >
            <Menu />
          </IconButton>
          <Typography variant="h6" noWrap>
            Persistent drawer
          </Typography>
        </Toolbar>
      </AppBar>
      <nav>
        <Drawer container={container} variant='temporary' open={openDrawer}>
          <div
            role="presentation"
            onClick={() => setDrawerOpen(!openDrawer)}
            onKeyDown={() => setDrawerOpen(!openDrawer)}
          >
            <div>
              <IconButton onClick={() => setDrawerOpen(!openDrawer)}>
                {<ChevronLeft />}
              </IconButton>
            </div>
            <Divider />
            <List>
              <ListItem button key={'Sort'}>
                <ListItemIcon>{<Sort/>}</ListItemIcon>
                <ListItemText primary={'Sort'} />
              </ListItem>
            </List>
          </div>
        </Drawer>
      </nav>
      <main style={{flexGrow: 1}}>
        <SortingPlugin algorithm={'bubble'}/>
      </main>
    </div>
  );
};

export default App;
